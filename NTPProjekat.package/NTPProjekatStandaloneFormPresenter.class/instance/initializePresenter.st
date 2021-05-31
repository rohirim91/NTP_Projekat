initialization
initializePresenter

	self submitSpButton action: [ self submitSp ].
	self submitPsoButton action: [ self submitPso ].
	self restoreButton action: [ self restore ]