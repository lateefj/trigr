local luatools = require "luatools"

test.is_source = function() 
	test.equal(luatools.is_source('file_example.lua'), true)
	test.equal(luatools.is_source('fileexample.go'), false)
end

test.is_test_source = function() 
	test.equal(luatools.is_test_source('file_example_test.lua'), true)
	test.equal(luatools.is_test_source('fileexample.lua'), false)
end
