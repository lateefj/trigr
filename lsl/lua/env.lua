-- Explicate exposure or no worky
local lsl = {}
lsl.env = {
  string = string,
  io = io,
  os = os,
  require = require,  -- XXX for testing
  module = module, -- XXX for testing
}

function lsl.run_with_env(env, fn, ...)
	for k, v in env_map() do
		env[k] = v
	end
	setfenv(fn, env)
	fn(...)
end

-- For running a string
function run_code_with_env(code)
	local fn = assert(loadstring(code))
	lsl.run_with_env(lsl.env, fn)
end

-- For running a file
function run_file_with_env(path)
	local file = assert(loadfile(path))
	lsl.run_with_env(lsl.env, file)
end

-- For running tests
function run_test_with_env(path, test_path, ...)
  test = utest()
  lsl.env.test = test
  lsl.env.test_path = test_path

  package.path = package.path .. ";" .. test_path .. "/?.lua"
  run_file_with_env(path, ...)
  -- Call for the results of the test
  local tests, failed =  test.result()
  -- Print out the summary of results
  test.summary()
end

return lsl
