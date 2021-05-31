specs
defaultSpec
	<spec: #default>

	| fontWidth formWidth |
	fontWidth := (StandardFonts defaultFont widthOfString: 'M').
	formWidth := fontWidth * 33.
	
	^ SpecLayout composed
		newRow: [ :row | 
			row add: #form right: formWidth.];
		yourself.	