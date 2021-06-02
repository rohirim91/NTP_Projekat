visualization
psoShowResults1Threshold

	| tab data b ds |
	tab := RTTabTable new input: self psoPosLocation asUrl retrieveContents usingDelimiter: $,.
	tab removeFirstRow.
	tab convertColumnsAsInteger: #(1).
	tab convertColumnsAsFloat: #(2 3).
	data := tab values.

	b := RTGrapher new.
	b extent: 750 @ 600.
	
	ds := RTData new.
	ds interaction fixedPopupText: [ :row | 't:' , (row value at: 1) asString , ' tsallis:' , (row value at: 2) asString ]. 
	ds dotShape ellipse 
		color: [:row | Color red:(0.7*(row value at: 2)**(1/4) * 255) + (0.3*(row value at: 3)*255) green:0 blue:255-(0.7*(row value at: 2)**(1/2.67) * 255) - (0.3*(row value at: 3)*255) alpha: 0.5 + (row value at: 3) ];
		size: [:row | 3 + ((row value at: 3) * 5)].
	ds points: data.
	
	ds x: [ :r | r value at: 1 ] .
	ds y: [ :r | r value at: 2 ].
	
	b add: ds. 

	b axisX noDecimal; title: 'Threshold value'.
	b axisY title: 'Tsallis entropy value'.
	b open 