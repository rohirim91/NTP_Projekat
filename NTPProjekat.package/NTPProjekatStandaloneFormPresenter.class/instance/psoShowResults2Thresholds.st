visualization
psoShowResults2Thresholds

	| tab data b ds |
	tab := RTTabTable new input: self psoPosLocation asUrl retrieveContents usingDelimiter: $,.
	tab removeFirstRow.
	tab convertColumnsAsInteger: #(1 2).
	tab convertColumnsAsFloat: #(3 4).
	data := tab values.

	b := RTGrapher new.
	b extent: 750 @ 600.
	
	ds := RTData new.
	ds interaction fixedPopupText: [ :row | 't1:' , (row value at: 1) asString , ' t2:' , (row value at: 2) asString , ' tsallis:' , (row value at: 3) asString ]. 
	ds dotShape ellipse 
		color: [:row | Color red:(0.7*(row value at: 3)**(1/4) * 255) + (0.3*(row value at: 4)*255) green:(0.4*(row value at: 3)**(1/2.67) * 255) + (0.6*(row value at: 4)*255) blue:255-(0.7*(row value at: 3)**(1/2.67) * 255) - (0.3*(row value at: 4)*255) alpha: 0.5 + (row value at: 4) ];
		size: [:row | 3 + ((row value at: 4) * 5) ].
	ds points: data.
	
	ds x: [ :r | r value at: 1 ] .
	ds y: [ :r | r value at: 2 ].
	
	b add: ds. 

	b axisX noDecimal; title: 'Threshold 1 value'.
	b axisY noDecimal; title: 'Threshold 2 value'.
	b open 