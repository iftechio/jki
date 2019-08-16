package image

import "strings"

// Image is docker image struct consist of domain, repo and tag
type Image struct {
	Domain string // domain or domain/namespace
	Repo   string
	Tag    string
}

func FromString(imageStr string) Image {
	image := Image{}
	parts := strings.Split(imageStr, "/")
	image.Domain = parts[0]
	repoWithTag := parts[len(parts)-1]
	parts = strings.Split(repoWithTag, ":")
	image.Repo = parts[0]
	if len(parts) == 1 {
		// missing colon
		image.Tag = "latest"
	} else {
		image.Tag = parts[1]
	}

	return image
}

func (image *Image) String() string {
	return image.Domain + "/" + image.Repo + ":" + image.Tag
}
