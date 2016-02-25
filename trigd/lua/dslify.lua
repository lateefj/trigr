local dsl = require 'lua.dsl'

local file = assert(loadfile(trig_dsl_path))

dsl.run_with_env(dsl.env, file)
