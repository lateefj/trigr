local dsl = {}

function dsl.run_with_env(env, fn, ...)
  setfenv(fn, env)
  fn(...)
end

-- Everything must be expose through this DSL
dsl.env = {
  string = string,
  io = io,
  os = os,
  print = print,
  console = function(fmt, ...)
    log.Info(string.format(fmt, ...))
  end,
  trig = trig,
  log = log,
}

return dsl
