package models

import (
	"mime/multipart"
)

type FileDownloadRequest struct {
	Namespace string `uri:"namespace" binding:"required"`
	Name      string `uri:"name" binding:"required"`
	Version   string `uri:"version"`
	Channel   string `uri:"channel"`
	FullNames []string
}

// IsLatest will return nil if we don't know its latest or not
func (f *FileDownloadRequest) IsLatest() *bool {
	var (
		isLatest  *bool = nil
		trueBool        = true
		falseBool       = false
	)
	if f.Channel != "" {
		isLatest = &falseBool
	}
	if f.Channel == "latest" {
		isLatest = &trueBool
	}

	// Note: adding this for backward compatibility
	// but ideally this shouldn't be the case
	if f.Version == "latest" {
		isLatest = &trueBool
	}
	return isLatest
}

// ToSnapshot creates snapshot
func (f *FileDownloadRequest) ToSnapshot() *Snapshot {
	s := &Snapshot{
		Namespace: f.Namespace,
		Name:      f.Name,
		Version:   f.Version,
	}
	if f.Channel == "latest" {
		s.Latest = true
	}

	// Note: adding this for backward compatibility
	// but ideally this shouldn't be the case
	if f.Version == "latest" {
		s.Latest = true
		s.Version = ""
	}
	return s
}

type DescriptorUploadRequest struct {
	Namespace string                `uri:"namespace" binding:"required"`
	Name      string                `form:"name" binding:"required"`
	Version   string                `form:"version" binding:"required,version"`
	File      *multipart.FileHeader `form:"file" binding:"required"`
	Latest    bool                  `form:"latest"`
	SkipRules []string              `form:"skiprules"`
	DryRun    bool                  `form:"dryrun"`
}

// ToSnapshot creates sanpshot
func (d *DescriptorUploadRequest) ToSnapshot() *Snapshot {
	return &Snapshot{
		Namespace: d.Namespace,
		Name:      d.Name,
		Version:   d.Version,
		Latest:    d.Latest,
	}
}
