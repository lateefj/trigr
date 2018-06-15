local gobuild = {}

gobuild.go_extension = ".go"

function gobuild.is_go_source(file_path)
  local extension = string.sub(file_path, #file_path-2, #file_path)
  if extension == gobuild.go_extension then
    print("returning true")
    return true
  end
  return false
end

function gobuild.run_tests(directory)
  print("Running go test in directory " .. directory)
  local t = io.popen("cd " .. directory .. "; go test")
  local output = t:read("*a")
  t:close()
  return output
end

return gobuild

