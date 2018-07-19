-- Create luatools package
local luatools = {}

-- Go extension
luatools.lua_extension = ".lua"
-- Check to see if file is a go source code
function luatools.is_source(file_path)
  -- Get the extension
  local extension = string.sub(file_path, #file_path-(#luatools.lua_extension -1), #file_path)
  -- Check to see if the extension matches
  if extension == luatools.lua_extension then
    return true
  end
  return false
end

-- Check to see if test is there
function luatools.is_test_source(file_path)
  if luatools.is_source(file_path) then
    local test_extension = "_test.lua"
    local extension = string.sub(file_path, #file_path-(#test_extension -1), #file_path)
    if test_extension == extension then
      return true
    end
  end
  return false
end

-- Run test for a directory
function luatools.run_test(file)
  --log("Running lua test file" .. file)
  return run_test_with_env(file)
end

-- Need this to build a package
return luatools
