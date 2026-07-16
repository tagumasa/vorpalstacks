package arn

// IoTBuilder provides methods for constructing AWS IoT ARNs.
type IoTBuilder struct{ *ARNBuilder }

// IoT returns an IoTBuilder for constructing IoT ARNs.
func (b *ARNBuilder) IoT() *IoTBuilder { return &IoTBuilder{b} }

func (b *IoTBuilder) Thing(name string) string {
	return b.Build("iot", "thing/"+name)
}

func (b *IoTBuilder) Certificate(id string) string {
	return b.Build("iot", "cert/"+id)
}

func (b *IoTBuilder) CACertificate(id string) string {
	return b.Build("iot", "cacert/"+id)
}

func (b *IoTBuilder) Policy(name string) string {
	return b.Build("iot", "policy/"+name)
}

func (b *IoTBuilder) Rule(name string) string {
	return b.Build("iot", "rule/"+name)
}

func (b *IoTBuilder) Job(id string) string {
	return b.Build("iot", "job/"+id)
}

func (b *IoTBuilder) ThingType(name string) string {
	return b.Build("iot", "thingtype/"+name)
}

func (b *IoTBuilder) ThingGroup(name string) string {
	return b.Build("iot", "thinggroup/"+name)
}

func (b *IoTBuilder) RoleAlias(name string) string {
	return b.Build("iot", "rolealias/"+name)
}

func (b *IoTBuilder) BillingGroup(name string) string {
	return b.Build("iot", "billinggroup/"+name)
}

func (b *IoTBuilder) Authorizer(name string) string {
	return b.Build("iot", "authorizer/"+name)
}

func (b *IoTBuilder) ProvisioningTemplate(name string) string {
	return b.Build("iot", "provisioningtemplate/"+name)
}

func (b *IoTBuilder) SecurityProfile(name string) string {
	return b.Build("iot", "securityprofile/"+name)
}

func (b *IoTBuilder) DomainConfiguration(name string) string {
	return b.Build("iot", "domainconfiguration/"+name)
}

func (b *IoTBuilder) CustomMetric(name string) string {
	return b.Build("iot", "custommetric/"+name)
}

func (b *IoTBuilder) Dimension(name string) string {
	return b.Build("iot", "dimension/"+name)
}

func (b *IoTBuilder) MitigationAction(name string) string {
	return b.Build("iot", "mitigationaction/"+name)
}

func (b *IoTBuilder) ScheduledAudit(name string) string {
	return b.Build("iot", "scheduledaudit/"+name)
}

func (b *IoTBuilder) FleetMetric(name string) string {
	return b.Build("iot", "fleetmetric/"+name)
}

func (b *IoTBuilder) TopicRuleDestination(name string) string {
	return b.Build("iot", "ruledestination/"+name)
}

func (b *IoTBuilder) JobTemplate(name string) string {
	return b.Build("iot", "jobtemplate/"+name)
}

func (b *IoTBuilder) Stream(name string) string {
	return b.Build("iot", "stream/"+name)
}

func (b *IoTBuilder) Package(name string) string {
	return b.Build("iot", "package/"+name)
}

func (b *IoTBuilder) Command(name string) string {
	return b.Build("iot", "command/"+name)
}

func (b *IoTBuilder) CertificateProvider(name string) string {
	return b.Build("iot", "certprovider/"+name)
}

func (b *IoTBuilder) OTAUpdate(name string) string {
	return b.Build("iot", "otaupdate/"+name)
}

// IoTEventsBuilder provides methods for constructing AWS IoT Events ARNs.
type IoTEventsBuilder struct{ *ARNBuilder }

// IoTEvents returns an IoTEventsBuilder for constructing IoT Events ARNs.
func (b *ARNBuilder) IoTEvents() *IoTEventsBuilder { return &IoTEventsBuilder{b} }

func (b *IoTEventsBuilder) DetectorModel(name string) string {
	return b.Build("iotevents", "detectorModel/"+name)
}

func (b *IoTEventsBuilder) Input(name string) string {
	return b.Build("iotevents", "input/"+name)
}

func (b *IoTEventsBuilder) AlarmModel(name string) string {
	return b.Build("iotevents", "alarmModel/"+name)
}
