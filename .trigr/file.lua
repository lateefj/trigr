-- Example of how to run the compiler based on 
-- Import go build module
local gobuild = require "gobuild"

local file_path = trig.Data["path"]

local filename = string.gsub(file_path, "(.*/)(.*)", "%2") 
-- Now the basepath
local basepath = string.sub(file_path, 0, #file_path - #filename)
-- If the extension is a go file then do custom commands
if gobuild.is_go_source(file_path) then
  -- Run test in directory
  local test_output = gobuild.run_tests(basepath)
  print(test_output)
  -- Run the build command
  [[--local m = io.popen("make build")
  print(m:read("*a"))
  m:close()
  -- Run the test command
  local test_make = io.popen("make test")
  print(test_make:read("*a"))
  test_make:close()--]]
  print("Done")
end
print(string.format("Type: %s", trig.Type))
print(string.format("Path: %s", trig.Data["path"]))
print(string.format("Operation: %s", trig.Data["op"]))

