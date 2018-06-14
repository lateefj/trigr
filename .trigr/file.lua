-- Example of how to run the compiler based on 
local file_path = trig.Data["path"]
local extension = string.sub(file_path, #file_path-2, #file_path)
-- If the extension is a go file then do custom commands
if extension == ".go" then
  -- Run the build command
  local m = io.popen("make build")
  print(m:read("*a"))
  m:close()
  -- Run the test command
  local test_make = io.popen("make test")
  print(test_make:read("*a"))
  test_make:close()
  print("Done")
end
print(string.format("Type: %s", trig.Type))
print(string.format("Path: %s", trig.Data["path"]))
print(string.format("Operation: %s", trig.Data["op"]))
