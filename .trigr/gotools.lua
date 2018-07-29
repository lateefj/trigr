-- Create gotools package
local gotools = {}

-- Go extension
gotools.go_extension = ".go"
-- Check to see if file is a go source code
function gotools.is_source(file_path)
  -- Get the extension
  local extension = string.sub(file_path, #file_path-2, #file_path)
  -- Check to see if the extension matches
  if extension == gotools.go_extension then
    return true
  end
  return false
end

-- Run test for a directory
function gotools.run_tests(directory)
  -- Go into the directory and run go test
  return exec("cd " .. directory .. "; go test")
end
-- Need this to build a package
return gotools

