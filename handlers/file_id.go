package handlers

import (
	"context"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func ListFilesInFolder(folderID string) ([]*drive.File, error) {

	ctx := context.Background()

	srv, err := drive.NewService(ctx, option.WithCredentialsFile("service.json"))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	query := fmt.Sprintf("'%s' in parents and mimeType contains 'image/'", folderID)
	call := srv.Files.List().Q(query).Fields("files(id, name, mimeType)")

	var files []*drive.File
	err = call.Pages(ctx, func(page *drive.FileList) error {
		files = append(files, page.Files...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve files: %v", err)
	}

	return files, nil
}

func GetFileMetadata(fileID string) (*drive.File, error) {
	ctx := context.Background()
	srv, err := drive.NewService(ctx, option.WithCredentialsFile("service.json"))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	file, err := srv.Files.Get(fileID).Fields("id, name, mimeType").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve file metadata: %v", err)
	}

	return file, nil
}
