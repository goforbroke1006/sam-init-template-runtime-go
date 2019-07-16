include ../../.env
export $(shell sed 's/=.*//' ../../.env)

debug:
	sam local invoke --event event.json HelloWorldFunction --debug
