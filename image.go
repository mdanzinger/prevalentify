package prevalentify

import "image"

// ImageSource represents the source of an image. This can be a URL, filepath, etc. The Image Resolver will resolve an image
// from the given ImageSource.
type ImageSource string

// Image represents a decoded image
type Image struct {
	Image  *image.Image
	Source ImageSource
}
