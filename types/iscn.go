package types

import (
	"time"
	
	iscntypes "github.com/likecoin/likechain/x/iscn/types"
)

// IscnParams represents the x/iscn parameters
type IscnParams struct {
	iscntypes.Params
	Height int64
}


// NewIscnParams allows to build a new IscnParams instance
func NewIscnParams(params iscntypes.Params, height int64) IscnParams {
	return IscnParams{
		Params: params,
		Height: height,
	}
}

type ContentMetadata struct {
	Context string
	Type string 
	About string
	Abstract string
	AccessMode string
	AcquireLicensePage string
	ArticleBody string
	Backstory string
	CopyrightHolder string
	CopyrightNotice string
	CopyrightYear string
	CreativeWorkStatus string
	Creator string
	DatePublished string
	Description string
	EncodingFormat string
	Headline string
	Keywords string
	License string
	LocationCreated string
	Name string
	Text string
	Url string
	UsageInfo string
	Version uint64
	WordCount string
}

type EntityStruct struct {
	Id string
	Name string
}

type Stakeholders struct {
	ContributionType string
	Entity EntityStruct
	RewardProportion uint64
}

type IscnRecord struct {
	Ipld string 
	Context string
	RecordID string
	RecordRoute string
	RecordType string
	ContentFingerprints []string
	ContentMetadata ContentMetadata
	RecordNotes string
	RecordTimestamp time.Time
	RecordVersion uint64
	Stakeholders []Stakeholders
	Height int64
}

// NewIscnRecord return a new IscnRecord instance
func NewIscnRecord(
	ipld string, 
	context string,
	recordID string,
	recordRoute string,
	recordType string,
	contentFingerprints []string,
	contentMetadata ContentMetadata,
	recordNotes string,
	recordTimestamp time.Time,
	recordVersion uint64,
	stakeholders []Stakeholders,
	height int64,
) IscnRecord {
	return IscnRecord{
		Ipld: ipld,
		Context: context,
		RecordID: recordID,
		RecordRoute: recordRoute,
		RecordType: recordType,
		ContentFingerprints: contentFingerprints,
		ContentMetadata: contentMetadata,
		RecordNotes: recordNotes,
		RecordTimestamp: recordTimestamp,
		RecordVersion: recordVersion,
		Stakeholders: stakeholders,
		Height: height,
	}
}

// Equal tells whether i IscnRecord and another IscnRecord contain the same data
func (i IscnRecord) Equal(another IscnRecord) bool {
	return i.Ipld == another.Ipld &&
		i.Context == another.Context &&
		i.RecordID == another.RecordID &&
		i.RecordRoute == another.RecordRoute &&
		i.RecordType == another.RecordType &&
		i.ContentFingerprints == another.ContentFingerprints &&
		i.ContentMetadata == another.ContentMetadata &&
		i.RecordNotes == another.RecordNotes &&
		i.RecordTimestamp.Equal(another.RecordTimestamp) &&
		i.RecordVersion == another.RecordVersion &&
		i.Stakeholders == another.Stakeholders &&
		i.Height == another.Height 
}