initialization
initializeWidgets

	spInputPathLabel := self newLabel label: 'Superpixels Input Path:'.
	spInputPathInput := self newTextInput autoAccept: true.

	spOutputPathLabel := self newLabel label: 'Superpixels Otput Path:'.
	spOutputPathInput := self newTextInput autoAccept: true.
	
	spParallelCheckboxLabel := self newLabel label: 'Parallel:'.
	spParallelCheckboxInput := self newCheckBox state: true.

	psoInputPathLabel := self newLabel label: 'PSO Input Path:'.
	psoInputPathInput := self newTextInput autoAccept: true.
	
	psoOutputPathLabel := self newLabel label: 'PSO Output Path:'.
	psoOutputPathInput := self newTextInput autoAccept: true.
	
	numParticlesLabel := self newLabel label: 'Num. Particles:'.
	numParticlesInput := self newTextInput autoAccept: true.

	maxIterLabel := self newLabel label: 'Max Iterations:'.
	maxIterInput := self newTextInput autoAccept: true.

	numThresholdsLabel := self newLabel label: 'Num. Thresholds:'.
	numThresholdsInput := self newTextInput autoAccept: true.
		
	psoParallelCheckboxLabel := self newLabel label: 'Parallel:'.
	psoParallelCheckboxInput := self newCheckBox state: true.
	
	submitSpButton := self newButton label: 'Superpixelize'.
	submitPsoButton := self newButton label: 'Apply PSO'.
	restoreButton := self newButton label: 'Restore'.

	self focusOrder
		add: spInputPathInput ;
		add: spOutputPathInput ;
		add: spParallelCheckboxInput ;
		add: psoInputPathInput ;
		add: psoOutputPathInput ;
		add: psoParallelCheckboxInput ;
		add: numParticlesInput ;
		add: maxIterInput ;
		add: numThresholdsInput .
		