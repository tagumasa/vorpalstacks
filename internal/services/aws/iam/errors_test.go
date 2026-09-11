package iam

import (
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
)

func TestIAMErrors(t *testing.T) {
	t.Run("predefined errors", func(t *testing.T) {
		assert.Equal(t, "LimitExceeded: Cannot exceed quota for AccessKeysPerUser: 2.", ErrAccessKeyLimitExceeded.Error())
		assert.Equal(t, http.StatusConflict, ErrAccessKeyLimitExceeded.GetHTTPStatusCode())

		assert.Equal(t, "PasswordPolicyViolation: The password does not meet the password policy requirements.", ErrPasswordPolicyViolation.Error())
		assert.Equal(t, http.StatusBadRequest, ErrPasswordPolicyViolation.GetHTTPStatusCode())

		assert.Equal(t, "MalformedPolicyDocument: This policy contains invalid JSON.", ErrMalformedPolicyDocument.Error())
		assert.Equal(t, http.StatusBadRequest, ErrMalformedPolicyDocument.GetHTTPStatusCode())

		assert.Equal(t, "LimitExceeded: Cannot exceed quota for PolicyVersions: 5.", ErrLimitExceededPolicyVersions.Error())
		assert.Equal(t, http.StatusConflict, ErrLimitExceededPolicyVersions.GetHTTPStatusCode())

		assert.Equal(t, "InvalidAuthenticationCode: Invalid authentication code.", ErrInvalidAuthenticationCode.Error())
		assert.Equal(t, http.StatusBadRequest, ErrInvalidAuthenticationCode.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchUserError", func(t *testing.T) {
		err := NewNoSuchUserError("testuser")
		assert.Equal(t, "NoSuchEntity: The user with name testuser cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewUserAlreadyExistsError", func(t *testing.T) {
		err := NewUserAlreadyExistsError("testuser")
		assert.Equal(t, "EntityAlreadyExists: User testuser already exists", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchAccessKeyError", func(t *testing.T) {
		err := NewNoSuchAccessKeyError("AKIAIOSFODNN7EXAMPLE")
		assert.Equal(t, "NoSuchEntity: The Access Key with name AKIAIOSFODNN7EXAMPLE cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchLoginProfileError", func(t *testing.T) {
		err := NewNoSuchLoginProfileError("testuser")
		assert.Equal(t, "NoSuchEntity: The Login profile for user with name testuser cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewLoginProfileAlreadyExistsError", func(t *testing.T) {
		err := NewLoginProfileAlreadyExistsError("testuser")
		assert.Equal(t, "EntityAlreadyExists: Login profile for user testuser already exists", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewDeleteConflictError", func(t *testing.T) {
		err := NewDeleteConflictError("must delete access keys first")
		assert.Equal(t, "DeleteConflict: must delete access keys first", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchGroupError", func(t *testing.T) {
		err := NewNoSuchGroupError("testgroup")
		assert.Equal(t, "NoSuchEntity: The group with name testgroup cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewGroupAlreadyExistsError", func(t *testing.T) {
		err := NewGroupAlreadyExistsError("testgroup")
		assert.Equal(t, "EntityAlreadyExists: Group testgroup already exists", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewUserNotInGroupError", func(t *testing.T) {
		err := NewUserNotInGroupError("testuser", "testgroup")
		assert.Equal(t, "NoSuchEntity: The User testuser in group with name testgroup cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchRoleError", func(t *testing.T) {
		err := NewNoSuchRoleError("testrole")
		assert.Equal(t, "NoSuchEntity: The role with name testrole cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewRoleAlreadyExistsError", func(t *testing.T) {
		err := NewRoleAlreadyExistsError("testrole")
		assert.Equal(t, "EntityAlreadyExists: Role testrole already exists", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewDeleteRoleConflictError", func(t *testing.T) {
		err := NewDeleteRoleConflictError("must detach all policies")
		assert.Equal(t, "DeleteConflict: must detach all policies", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchInstanceProfileError", func(t *testing.T) {
		err := NewNoSuchInstanceProfileError("testprofile")
		assert.Equal(t, "NoSuchEntity: The Instance Profile with name testprofile cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewInstanceProfileAlreadyExistsError", func(t *testing.T) {
		err := NewInstanceProfileAlreadyExistsError("testprofile")
		assert.Equal(t, "EntityAlreadyExists: Instance Profile testprofile already exists", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewRoleNotInInstanceProfileError", func(t *testing.T) {
		err := NewRoleNotInInstanceProfileError("testrole", "testprofile")
		assert.Equal(t, "NoSuchEntity: The Role testrole in instance profile with name testprofile cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewRoleAlreadyInInstanceProfileError", func(t *testing.T) {
		err := NewRoleAlreadyInInstanceProfileError("testrole", "testprofile")
		assert.Equal(t, "LimitExceededException: Cannot exceed quota for InstanceProfilesPerRole: 1. Already associated with instance profile testprofile.", err.Error())
		assert.Equal(t, http.StatusBadRequest, err.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchPolicyError", func(t *testing.T) {
		err := NewNoSuchPolicyError("arn:aws:iam::123456789012:policy/test")
		assert.Equal(t, "NoSuchEntity: The policy with name arn:aws:iam::123456789012:policy/test cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewPolicyAlreadyExistsError", func(t *testing.T) {
		err := NewPolicyAlreadyExistsError("testpolicy")
		assert.Equal(t, "EntityAlreadyExists: Policy testpolicy already exists", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewDeletePolicyConflictError", func(t *testing.T) {
		err := NewDeletePolicyConflictError("arn:aws:iam::123456789012:policy/test")
		assert.Equal(t, "DeleteConflict: Cannot delete policy arn:aws:iam::123456789012:policy/test, there are attachments.", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchPolicyVersionError", func(t *testing.T) {
		err := NewNoSuchPolicyVersionError("v1")
		assert.Equal(t, "NoSuchEntity: The Policy version with name v1 cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchMFADeviceError", func(t *testing.T) {
		err := NewNoSuchMFADeviceError("arn:aws:iam::123456789012:mfa/testuser")
		assert.Equal(t, "NoSuchEntity: The MFA Device with name arn:aws:iam::123456789012:mfa/testuser cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewMFADeviceAlreadyAssignedError", func(t *testing.T) {
		err := NewMFADeviceAlreadyAssignedError("arn:aws:iam::123456789012:mfa/testuser")
		assert.Equal(t, "EntityAlreadyExists: MFA Device arn:aws:iam::123456789012:mfa/testuser already exists", err.Error())
		assert.Equal(t, http.StatusConflict, err.GetHTTPStatusCode())
	})

	t.Run("NewNoSuchPasswordPolicyError", func(t *testing.T) {
		err := NewNoSuchPasswordPolicyError()
		assert.Equal(t, "NoSuchEntity: The Password Policy with name  cannot be found.", err.Error())
		assert.Equal(t, http.StatusNotFound, err.GetHTTPStatusCode())
	})

	t.Run("NewInvalidInputError", func(t *testing.T) {
		err := NewInvalidInputError("RoleName", "invalid characters")
		assert.Equal(t, "InvalidInput: The input parameter RoleName is invalid: invalid characters", err.Error())
		assert.Equal(t, http.StatusBadRequest, err.GetHTTPStatusCode())
	})
}

// resolveUserName interpolates the caller's principal when an omitted
// UserName cannot be defaulted — no placeholder text reaches the wire.
func TestResolveUserNameInterpolatesPrincipal(t *testing.T) {
	_, err := resolveUserName(&request.RequestContext{
		PrincipalType: request.PrincipalTypeRole,
		Principal:     "arn:aws:sts::123456789012:assumed-role/deploy/session",
	}, "")
	if err == nil {
		t.Fatal("a non-user caller with no UserName must be rejected")
	}
	awsErr, ok := err.(*awserrors.AWSError)
	if !ok {
		t.Fatalf("expected *awserrors.AWSError, got %T", err)
	}
	assert.Equal(t, http.StatusNotFound, awsErr.GetHTTPStatusCode())
	assert.Contains(t, awsErr.Error(), "arn:aws:sts::123456789012:assumed-role/deploy/session")
	assert.NotContains(t, awsErr.Error(), "{")
}

// The empty-required-parameter paths are InvalidInput-class validation
// errors naming the missing member (each swapped member is required in the
// Smithy model) — a wire-visible class change from the former NoSuchEntity
// sentinels, pinned per site. The SDK's client-side validation makes these
// server-side paths unreachable through the AWS SDK, so they are pinned as
// unit tests.
func TestRequiredParameterErrorsAreValidationClass(t *testing.T) {
	validationErr := func(t *testing.T, err error, param string) {
		t.Helper()
		awsErr, ok := err.(*awserrors.AWSError)
		if !ok {
			t.Fatalf("%s: expected *awserrors.AWSError, got %T", param, err)
		}
		assert.Equal(t, http.StatusBadRequest, awsErr.GetHTTPStatusCode(), param)
		assert.Equal(t, "InvalidInput: Required parameter "+param+" is missing.", awsErr.Error())
	}

	t.Run("principalNameRequiredError", func(t *testing.T) {
		validationErr(t, principalNameRequiredError(PrincipalTypeUser), "UserName")
		validationErr(t, principalNameRequiredError(PrincipalTypeGroup), "GroupName")
		validationErr(t, principalNameRequiredError(PrincipalTypeRole), "RoleName")
		validationErr(t, principalNameRequiredError("openid-provider"), "PrincipalName")
	})

	// The empty-name branch returns before any store access, so a nil store
	// keeps these calls hermetic.
	t.Run("deleteGroupCore", func(t *testing.T) {
		s := &IAMService{}
		validationErr(t, s.deleteGroupCore(nil, &DeleteGroupInput{}), "GroupName")
	})

	t.Run("deleteRoleCore", func(t *testing.T) {
		s := &IAMService{}
		validationErr(t, s.deleteRoleCore(nil, &DeleteRoleInput{}), "RoleName")
	})

	t.Run("createVirtualMFADeviceCore", func(t *testing.T) {
		s := &IAMService{}
		_, err := s.createVirtualMFADeviceCore(nil, &CreateVirtualMFADeviceInput{})
		validationErr(t, err, "VirtualMFADeviceName")
	})
}

// TestAWSErrorMessagesCarryNoPlaceholders is the gate for the single
// error-construction system: no error-constructor call in the package may
// embed a curly-brace placeholder in a string argument. The error message
// is the final wire text — an unsubstituted "{Placeholder}" would reach
// the client verbatim.
func TestAWSErrorMessagesCarryNoPlaceholders(t *testing.T) {
	constructorName := regexp.MustCompile(`^New\w*(Error|Exception)$`)
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read package directory: %v", err)
	}
	fset := token.NewFileSet()
	checked := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, entry.Name(), nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", entry.Name(), err)
		}
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			var name string
			switch fn := call.Fun.(type) {
			case *ast.Ident:
				name = fn.Name
			case *ast.SelectorExpr:
				name = fn.Sel.Name
			default:
				return true
			}
			if !constructorName.MatchString(name) {
				return true
			}
			checked++
			for _, arg := range call.Args {
				ast.Inspect(arg, func(m ast.Node) bool {
					lit, ok := m.(*ast.BasicLit)
					if !ok || lit.Kind != token.STRING {
						return true
					}
					if strings.Contains(lit.Value, "{") {
						t.Errorf("%s: %s argument %s embeds a placeholder", fset.Position(call.Pos()), name, lit.Value)
					}
					return true
				})
			}
			return true
		})
	}
	if checked == 0 {
		t.Fatal("no error-constructor calls found — the gate is not scanning the package")
	}
}
