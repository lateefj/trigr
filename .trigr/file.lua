-- Example of how to run the compiler based on 
-- Import go build module
local gotools = require "gotools"
-- Import lua tools
local luatools = require "luatools"

-- Get the path from the trigger event data
local file_path = trig.Data["path"]

-- Get the filename
local filename = string.gsub(file_path, "(.*/)(.*)", "%2") 
-- Now the basepath
local basepath = string.sub(file_path, 0, #file_path - #filename)
-- List of supported file operations
local supported_ops = { 'write', 'create', 'remove', 'rename' }
--log("Path is " .. file_path)
-- Make sure it is a supported op
if contains(supported_ops, trig.Data["op"]) then
  -- If the extension is a go file then do custom commands
  if gotools.is_source(file_path) then
    -- Run test in directory
    log(gotools.run_tests(exec, basepath))
    -- TODO: Run the build command
  end

  -- If the extension is a lua file then do custom commands
  if luatools.is_test_source(file_path) then
    -- Run test in directory
    luatools.run_test(file_path)
  end
end 
--[[print(string.format("Trig Type: %s", trig.Type))
print("Trig Data key : value")
for k, v in trig.Data() do
  print(string.format("\t%s : %v", k, v))
end
]]--

