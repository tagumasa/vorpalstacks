package testutil

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func (r *TestRunner) runRoute53CidrCollectionTests(tc *route53TestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("route53", "CidrCollection_Lifecycle", func() error {
		name := tc.uniqName("cidr-life")
		createResp, err := tc.createCidrCollection(name, tc.callerRef("cidrlife"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		if createResp.Collection == nil || aws.ToString(createResp.Collection.Id) == "" {
			return fmt.Errorf("collection or id is nil")
		}
		if aws.ToString(createResp.Collection.Name) != name {
			return fmt.Errorf("name mismatch: got %q", aws.ToString(createResp.Collection.Name))
		}
		if !strings.Contains(aws.ToString(createResp.Collection.Arn), aws.ToString(createResp.Collection.Id)) {
			return fmt.Errorf("arn does not embed the collection id: %q", aws.ToString(createResp.Collection.Arn))
		}
		collectionID := aws.ToString(createResp.Collection.Id)

		defer tc.deleteCidrCollection(collectionID)

		changeResp, err := tc.client.ChangeCidrCollection(tc.ctx, &route53.ChangeCidrCollectionInput{
			Id:                aws.String(collectionID),
			CollectionVersion: aws.Int64(1),
			Changes: []types.CidrCollectionChange{
				{
					LocationName: aws.String("alpha-loc"),
					Action:       types.CidrCollectionChangeActionPut,
					CidrList:     []string{"10.150.0.0/24"},
				},
				{
					LocationName: aws.String("beta-loc"),
					Action:       types.CidrCollectionChangeActionPut,
					CidrList:     []string{"10.151.0.0/24"},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("change PUT: %v", err)
		}
		// The response Id is a change id: GetChange must resolve it and
		// report the synchronously applied batch as INSYNC.
		changeID := aws.ToString(changeResp.Id)
		if changeID == "" {
			return fmt.Errorf("change response id is empty")
		}
		getChangeResp, err := tc.client.GetChange(tc.ctx, &route53.GetChangeInput{Id: aws.String(changeID)})
		if err != nil {
			return fmt.Errorf("get change %q: %v", changeID, err)
		}
		if getChangeResp.ChangeInfo == nil || getChangeResp.ChangeInfo.Status != types.ChangeStatusInsync {
			return fmt.Errorf("expected INSYNC change info, got %+v", getChangeResp.ChangeInfo)
		}

		locResp, err := tc.client.ListCidrLocations(tc.ctx, &route53.ListCidrLocationsInput{
			CollectionId: aws.String(collectionID),
		})
		if err != nil {
			return fmt.Errorf("list locations: %v", err)
		}
		found := map[string]bool{}
		for _, loc := range locResp.CidrLocations {
			found[aws.ToString(loc.LocationName)] = true
		}
		if !found["alpha-loc"] || !found["beta-loc"] {
			return fmt.Errorf("locations alpha-loc/beta-loc not listed: %v", found)
		}

		// The location filter must narrow to that location's blocks only.
		blockResp, err := tc.client.ListCidrBlocks(tc.ctx, &route53.ListCidrBlocksInput{
			CollectionId: aws.String(collectionID),
			LocationName: aws.String("alpha-loc"),
		})
		if err != nil {
			return fmt.Errorf("list blocks: %v", err)
		}
		if len(blockResp.CidrBlocks) != 1 || aws.ToString(blockResp.CidrBlocks[0].CidrBlock) != "10.150.0.0/24" {
			return fmt.Errorf("expected exactly 1 block 10.150.0.0/24 for alpha-loc, got %v", blockResp.CidrBlocks)
		}

		// A location that matches nothing is rejected, not emptied.
		_, err = tc.client.ListCidrBlocks(tc.ctx, &route53.ListCidrBlocksInput{
			CollectionId: aws.String(collectionID),
			LocationName: aws.String("no-such-loc"),
		})
		return expectRoute53Error(err, "NoSuchCidrLocationException", 404)
	}))

	results = append(results, r.RunTest("route53", "CidrCollection_Pagination", func() error {
		createResp, err := tc.createCidrCollection(tc.uniqName("cidr-page"), tc.callerRef("cidrpage"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		collectionID := aws.ToString(createResp.Collection.Id)

		defer tc.deleteCidrCollection(collectionID)

		changes := []types.CidrCollectionChange{
			{LocationName: aws.String("page-a"), Action: types.CidrCollectionChangeActionPut, CidrList: []string{"10.152.1.0/24"}},
			{LocationName: aws.String("page-b"), Action: types.CidrCollectionChangeActionPut, CidrList: []string{"10.152.2.0/24"}},
			{LocationName: aws.String("page-c"), Action: types.CidrCollectionChangeActionPut, CidrList: []string{"10.152.3.0/24"}},
		}
		if _, err := tc.client.ChangeCidrCollection(tc.ctx, &route53.ChangeCidrCollectionInput{
			Id:                aws.String(collectionID),
			CollectionVersion: aws.Int64(1),
			Changes:           changes,
		}); err != nil {
			return fmt.Errorf("change PUT: %v", err)
		}

		// Location pages: 2 + 1 with the continuation token.
		page1, err := tc.client.ListCidrLocations(tc.ctx, &route53.ListCidrLocationsInput{
			CollectionId: aws.String(collectionID),
			MaxResults:   aws.Int32(2),
		})
		if err != nil {
			return fmt.Errorf("list locations page 1: %v", err)
		}
		if len(page1.CidrLocations) != 2 {
			return fmt.Errorf("expected 2 locations on page 1, got %d", len(page1.CidrLocations))
		}
		if aws.ToString(page1.NextToken) == "" {
			return fmt.Errorf("expected NextToken on truncated location page")
		}
		page2, err := tc.client.ListCidrLocations(tc.ctx, &route53.ListCidrLocationsInput{
			CollectionId: aws.String(collectionID),
			NextToken:    page1.NextToken,
		})
		if err != nil {
			return fmt.Errorf("list locations page 2: %v", err)
		}
		if len(page2.CidrLocations) != 1 || aws.ToString(page2.NextToken) != "" {
			return fmt.Errorf("expected 1 location and no token on page 2, got %d locations, token %q",
				len(page2.CidrLocations), aws.ToString(page2.NextToken))
		}

		// Block pages: walk with MaxResults=2 and collect every cidr.
		seenBlocks := map[string]bool{}
		var blockToken *string
		pages := 0
		for {
			resp, err := tc.client.ListCidrBlocks(tc.ctx, &route53.ListCidrBlocksInput{
				CollectionId: aws.String(collectionID),
				MaxResults:   aws.Int32(2),
				NextToken:    blockToken,
			})
			if err != nil {
				return fmt.Errorf("list blocks page %d: %v", pages+1, err)
			}
			for _, b := range resp.CidrBlocks {
				seenBlocks[aws.ToString(b.CidrBlock)] = true
			}
			pages++
			if aws.ToString(resp.NextToken) == "" || pages > 5 {
				break
			}
			blockToken = resp.NextToken
		}
		if pages != 2 || len(seenBlocks) != 3 {
			return fmt.Errorf("expected 2 block pages covering 3 cidrs, got %d pages, %d cidrs", pages, len(seenBlocks))
		}

		// Blocks whose insertion order differs from lexical order must
		// survive pagination: 10.152.9.0/24 is PUT before 10.152.10.0/24,
		// which sorts ahead of it, so a cursor walk with one block per
		// page must still reach both.
		if _, err := tc.client.ChangeCidrCollection(tc.ctx, &route53.ChangeCidrCollectionInput{
			Id:                aws.String(collectionID),
			CollectionVersion: aws.Int64(2),
			Changes: []types.CidrCollectionChange{{
				LocationName: aws.String("page-multi"),
				Action:       types.CidrCollectionChangeActionPut,
				CidrList:     []string{"10.152.9.0/24", "10.152.10.0/24"},
			}},
		}); err != nil {
			return fmt.Errorf("change PUT page-multi: %v", err)
		}
		gotMulti := map[string]bool{}
		multiPages := 0
		var multiToken *string
		for {
			resp, err := tc.client.ListCidrBlocks(tc.ctx, &route53.ListCidrBlocksInput{
				CollectionId: aws.String(collectionID),
				LocationName: aws.String("page-multi"),
				MaxResults:   aws.Int32(1),
				NextToken:    multiToken,
			})
			if err != nil {
				return fmt.Errorf("list page-multi blocks page %d: %v", multiPages+1, err)
			}
			for _, b := range resp.CidrBlocks {
				gotMulti[aws.ToString(b.CidrBlock)] = true
			}
			multiPages++
			if aws.ToString(resp.NextToken) == "" || multiPages > 4 {
				break
			}
			multiToken = resp.NextToken
		}
		if multiPages != 2 || !gotMulti["10.152.9.0/24"] || !gotMulti["10.152.10.0/24"] {
			return fmt.Errorf("expected 2 single-block pages covering both page-multi cidrs, got %d pages, %v",
				multiPages, gotMulti)
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "CidrCollection_MalformedToken", func() error {
		createResp, err := tc.createCidrCollection(tc.uniqName("cidr-token"), tc.callerRef("cidrtoken"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		collectionID := aws.ToString(createResp.Collection.Id)

		defer tc.deleteCidrCollection(collectionID)

		if _, err := tc.client.ChangeCidrCollection(tc.ctx, &route53.ChangeCidrCollectionInput{
			Id:                aws.String(collectionID),
			CollectionVersion: aws.Int64(1),
			Changes: []types.CidrCollectionChange{{
				LocationName: aws.String("token-loc"),
				Action:       types.CidrCollectionChangeActionPut,
				CidrList:     []string{"10.154.1.0/24"},
			}},
		}); err != nil {
			return fmt.Errorf("change PUT: %v", err)
		}

		// A continuation token that cannot be decoded is rejected instead
		// of silently restarting the listing from the first page.
		_, err = tc.client.ListCidrBlocks(tc.ctx, &route53.ListCidrBlocksInput{
			CollectionId: aws.String(collectionID),
			NextToken:    aws.String("!!!not-a-valid-token!!!"),
		})
		if err := expectRoute53Error(err, "InvalidInput", 400); err != nil {
			return err
		}

		// Valid base64 without the location/cidr cursor split is equally
		// invalid for the blocks listing.
		sepless := base64.StdEncoding.EncodeToString([]byte("token-loc"))
		_, err = tc.client.ListCidrBlocks(tc.ctx, &route53.ListCidrBlocksInput{
			CollectionId: aws.String(collectionID),
			NextToken:    aws.String(sepless),
		})
		if err := expectRoute53Error(err, "InvalidInput", 400); err != nil {
			return err
		}

		_, err = tc.client.ListCidrLocations(tc.ctx, &route53.ListCidrLocationsInput{
			CollectionId: aws.String(collectionID),
			NextToken:    aws.String("!!!not-a-valid-token!!!"),
		})
		if err := expectRoute53Error(err, "InvalidInput", 400); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "ChangeCidrCollection_TooManyChanges", func() error {
		createResp, err := tc.createCidrCollection(tc.uniqName("cidr-cap"), tc.callerRef("cidrcap"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		collectionID := aws.ToString(createResp.Collection.Id)

		defer tc.deleteCidrCollection(collectionID)

		// The Changes list is capped at 1000 entries per request.
		changes := make([]types.CidrCollectionChange, 1001)
		for i := range changes {
			changes[i] = types.CidrCollectionChange{
				LocationName: aws.String(fmt.Sprintf("cap-%04d", i)),
				Action:       types.CidrCollectionChangeActionPut,
				CidrList:     []string{"10.153.0.0/24"},
			}
		}
		_, err = tc.client.ChangeCidrCollection(tc.ctx, &route53.ChangeCidrCollectionInput{
			Id:                aws.String(collectionID),
			CollectionVersion: aws.Int64(1),
			Changes:           changes,
		})
		return expectRoute53Error(err, "InvalidInput", 400)
	}))

	results = append(results, r.RunTest("route53", "CreateCidrCollection_DuplicateName", func() error {
		name := tc.uniqName("cidr-dup")
		createResp, err := tc.createCidrCollection(name, tc.callerRef("cidrdup-a"))
		if err != nil {
			return fmt.Errorf("create first: %v", err)
		}
		collectionID := aws.ToString(createResp.Collection.Id)

		defer tc.deleteCidrCollection(collectionID)

		// The same name with a different client token is rejected.
		_, err = tc.createCidrCollection(name, tc.callerRef("cidrdup-b"))
		if err := expectRoute53Error(err, "CidrCollectionAlreadyExistsException", 400); err != nil {
			return err
		}

		// A retry with the original client token returns the existing
		// collection instead of creating a second one.
		retryResp, err := tc.createCidrCollection(name, tc.callerRef("cidrdup-a"))
		if err != nil {
			return fmt.Errorf("same-caller-reference retry: %v", err)
		}
		if aws.ToString(retryResp.Collection.Id) != collectionID {
			return fmt.Errorf("retry returned a different collection: got %q want %q",
				aws.ToString(retryResp.Collection.Id), collectionID)
		}
		return nil
	}))

	results = append(results, r.RunTest("route53", "ChangeCidrCollection_VersionMismatch", func() error {
		createResp, err := tc.createCidrCollection(tc.uniqName("cidr-ver"), tc.callerRef("cidrver"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		collectionID := aws.ToString(createResp.Collection.Id)

		defer tc.deleteCidrCollection(collectionID)

		_, err = tc.client.ChangeCidrCollection(tc.ctx, &route53.ChangeCidrCollectionInput{
			Id:                aws.String(collectionID),
			CollectionVersion: aws.Int64(999),
			Changes: []types.CidrCollectionChange{{
				LocationName: aws.String("loc"),
				Action:       types.CidrCollectionChangeActionPut,
				CidrList:     []string{"10.151.0.0/24"},
			}},
		})
		return expectRoute53Error(err, "CidrCollectionVersionMismatchException", 409)
	}))

	results = append(results, r.RunTest("route53", "ListCidrLocations_NonExistentCollection", func() error {
		_, err := tc.client.ListCidrLocations(tc.ctx, &route53.ListCidrLocationsInput{
			CollectionId: aws.String("9999999999999999999"),
		})
		return expectRoute53Error(err, "NoSuchCidrCollectionException", 404)
	}))

	// The change member constraints from the API model: a location name is
	// 1-16 characters from [0-9A-Za-z_-], a CIDR entry is 1-50 non-blank
	// characters, and each change's Cidr list holds at most 1000 entries.
	results = append(results, r.RunTest("route53", "ChangeCidrCollection_MemberLimitsRejected", func() error {
		createResp, err := tc.createCidrCollection(tc.uniqName("cidr-member"), tc.callerRef("cidrmember"))
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		collectionID := aws.ToString(createResp.Collection.Id)

		defer tc.deleteCidrCollection(collectionID)

		rejections := []types.CidrCollectionChange{
			{LocationName: aws.String("overlong-location"), Action: types.CidrCollectionChangeActionPut, CidrList: []string{"10.150.0.0/24"}},
			{LocationName: aws.String("bad loc"), Action: types.CidrCollectionChangeActionPut, CidrList: []string{"10.150.0.0/24"}},
			{LocationName: aws.String("valid-loc"), Action: types.CidrCollectionChangeActionPut, CidrList: []string{"   "}},
		}
		for i, change := range rejections {
			_, err := tc.client.ChangeCidrCollection(tc.ctx, &route53.ChangeCidrCollectionInput{
				Id:                aws.String(collectionID),
				CollectionVersion: aws.Int64(1),
				Changes:           []types.CidrCollectionChange{change},
			})
			if err := expectRoute53Error(err, "InvalidInput", 400); err != nil {
				return fmt.Errorf("change %d: %w", i, err)
			}
		}

		overLengthList := make([]string, 1001)
		for i := range overLengthList {
			overLengthList[i] = fmt.Sprintf("10.160.%d.%d/24", i/256, i%256)
		}
		_, err = tc.client.ChangeCidrCollection(tc.ctx, &route53.ChangeCidrCollectionInput{
			Id:                aws.String(collectionID),
			CollectionVersion: aws.Int64(1),
			Changes: []types.CidrCollectionChange{{
				LocationName: aws.String("list-cap"),
				Action:       types.CidrCollectionChangeActionPut,
				CidrList:     overLengthList,
			}},
		})
		return expectRoute53Error(err, "InvalidInput", 400)
	}))

	return results
}
