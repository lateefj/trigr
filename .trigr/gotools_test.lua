
local gotools = require "gotools"
-- Test that the file source matching code works
test.is_source = function() 
	test.equal(gotools.is_source('fileex.go'), true)
	test.equal(gotools.is_source('fileexample.lua'), false)
end

