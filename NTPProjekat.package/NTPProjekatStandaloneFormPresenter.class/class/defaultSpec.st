specs
defaultSpec
	<spec: #default>

	| fontWidth labelsWidth rowHeight |
	fontWidth := (StandardFonts defaultFont widthOfString: 'M').
	labelsWidth := fontWidth * 12.
	rowHeight := self inputTextHeight.
	
	^ SpecLayout composed
		newColumn: [ :col | 

			col newRow: [ :row | 
				row add: #spInputPathLabel right: labelsWidth.
				row add: #spInputPathInput ]
			height: rowHeight.
		
			col newRow: [ :row | 
				row add: #spOutputPathLabel right: labelsWidth.
				row add: #spOutputPathInput ]
			height: rowHeight.
			
			col newRow: [ :row | 
				row add: #spParallelCheckboxLabel right: labelsWidth.
				row add: #spParallelCheckboxInput ]
			height: rowHeight.

			col newRow: [ :row | 
				row add: #psoInputPathLabel right: labelsWidth.
				row add: #psoInputPathInput ]
			height: rowHeight.

			col newRow: [ :row | 
				row add: #psoOutputPathLabel right: labelsWidth.
				row add: #psoOutputPathInput ]
			height: rowHeight.
			
		col newRow: [ :row | 
				row add: #psoParallelCheckboxLabel right: labelsWidth.
				row add: #psoParallelCheckboxInput ]
			height: rowHeight.

		col newRow: [ :row | 
				row add: #numParticlesLabel right: labelsWidth.
				row add: #numParticlesInput ]
			height: rowHeight.

		col newRow: [ :row | 
				row add: #maxIterLabel right: labelsWidth.
				row add: #maxIterInput ]
			height: rowHeight.

		col newRow: [ :row | 
				row add: #numThresholdsLabel right: labelsWidth.
				row add: #numThresholdsInput ]
			height: rowHeight.
	
	col newRow: [ :row | 
				row add: #submitSpButton.
				row add: #submitPsoButton.
				row add: #restoreButton ]
			height: rowHeight*1.3.

		col newRow: [ :row |  ] ];
			
		yourself.	