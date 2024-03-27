majority:
	@mkdir -p results/majority/1000-500-${threshold}
	@cd results/majority/1000-500-${threshold} && ./../../../bin/majority-cellular-automata 1000 500 ${threshold}

proximity:
	@mkdir -p results/proximity/1000-500
	@cd results/proximity/1000-500 && ./../../../bin/proximity-cellular-automata 1000 500
