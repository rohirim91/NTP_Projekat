initialization
initializeWidgets
	form := (self instantiate: SpecDemoLabeledContainer)
		label: 'Algorithm Settings';
		subwidget: (self instantiate: NTPProjekatStandaloneFormPresenter  on: self announcingObject).
	self focusOrder
		add: form.