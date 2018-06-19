-- Example of how to run the compiler based on 
-- Import go build module
local gobuild = require "gobuild"

-- Get the path from the trigger event data
local file_path = trig.Data["path"]

-- Get the filename
local filename = string.gsub(file_path, "(.*/)(.*)", "%2") 
-- Now the basepath
local basepath = string.sub(file_path, 0, #file_path - #filename)
-- List of supported file operations
local supported_ops = { 'write', 'create', 'remove' }

-- If the extension is a go file then do custom commands
if contains(supported_ops, trig.Data["op"]) and gobuild.is_go_source(file_path) then
  -- Run test in directory
  print(gobuild.run_tests(basepath))
  -- Run the build command
  print("Done handling go source")
end
print(string.format("Type: %s", trig.Type))
print(string.format("Path: %s", trig.Data["path"]))
print(string.format("Operation: %s", trig.Data["op"]))

