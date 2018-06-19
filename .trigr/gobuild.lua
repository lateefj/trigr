-- Create gobuild package
local gobuild = {}

-- Go extension
gobuild.go_extension = ".go"
-- Check to see if file is a go source code
function gobuild.is_go_source(file_path)
  -- Get the extension
  local extension = string.sub(file_path, #file_path-2, #file_path)
  -- Check to see if the extension matches
  if extension == gobuild.go_extension then
    return true
  end
  return false
end

-- Run test for a directory
function gobuild.run_tests(directory)
  print("Running go test in directory " .. directory)
  -- Go into the directory and run go test
  local t = io.popen("cd " .. directory .. "; go test")
  -- Store output into a variable
  local output = t:read("*a")
  -- Close the connection
  t:close()
  return output
end
-- Need this to build a package
return gobuild

