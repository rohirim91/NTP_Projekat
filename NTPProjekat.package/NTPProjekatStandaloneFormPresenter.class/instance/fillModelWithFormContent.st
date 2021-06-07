model updates
fillModelWithFormContent

	self workingModel
		spInputPath: self spInputPathInput text;
		spOutputPath: self spOutputPathInput text;
		spParallel: self spParallelCheckboxInput state;
		psoInputPath: self psoInputPathInput text;
		psoOutputPath: self psoOutputPathInput text;
		psoParallel: self psoParallelCheckboxInput state;
		numParticles: numParticlesInput text asNumber;
		maxIter: maxIterInput text asNumber;
		numThresholds: numThresholdsInput text asNumber.