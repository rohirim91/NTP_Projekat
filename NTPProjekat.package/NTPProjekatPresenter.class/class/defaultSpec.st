specs
defaultSpec
	<spec: #default>

	| fontWidth formWidth |
	fontWidth := (StandardFonts defaultFont widthOfString: 'M').
	formWidth := fontWidth * 33.
		
	^ SpecRowLayout composed
		newRow: [ :row | 
			row add: #form right: formWidth .] ;
		yourself.	