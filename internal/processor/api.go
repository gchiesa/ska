package processor

import "os"

func (tp *FileTreeProcessor) Render(withVariables map[string]interface{}) error {
	if tp.workingDir == "" {
		var err error
		tp.workingDir, err = os.MkdirTemp("", "swansonRenderer")
		if err != nil {
			return err
		}
	}

	// WAVE 1 - render the tree structure
	// allocate the folders and files in a staging directory to be rendered
	if err := tp.buildStagingFileTree(withVariables); err != nil {
		return err
	}

	// WAVE 2 - decompose the swanson managed partials
	// create a set of partials that are related to the files in the staging directory
	if err := tp.loadMultiparts(); err != nil {
		return err
	}

	// WAVE 3 - expand template
	// render all the templates, but if a partial exists for a file then expands only the partials
	if err := tp.renderStagingFileTree(withVariables); err != nil {
		return err
	}

	// WAVE 4 - copy to destination the staging directory
	// copy the staging directory to the destination with the process
	// for each file (non-swanson) copy the file first
	// then replace the partials with the expanded content
	// **IF the file mustBeSkipped then skip, otherwise copy
	// **IF the file already exists in the destination then
	// only replace the partials with the expanded content
	if err := tp.updateDestinationFileTree(); err != nil {
		return err
	}
	return nil
}
