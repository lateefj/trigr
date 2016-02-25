local dsl = {}
if not setfenv then -- Lua 5.2
  -- based on http://lua-users.org/lists/lua-l/2010-06/msg00314.html
  -- this assumes f is a function
  local function findenv(f)
    local level = 1
    repeat
      local name, value = debug.getupvalue(f, level)
      if name == '_ENV' then return level, value end
      level = level + 1
    until name == nil
    return nil end
  getfenv = function (f) return(select(2, findenv(f)) or _G) end
  setfenv = function (f, t)
    local level = findenv(f)
    if level then debug.setupvalue(f, level, t) end
    return f end
end


function dsl.run_with_env(env, fn, ...)
  setfenv(fn, env)
  fn(...)
end


dsl.env = {
  move = function(ops)
    print("I moved to", ops[1], ops[2])
  end,
  speak = function(message)
    print("I said:", message)
  end,
  print = function(...)
    print(...)
  end,
  console = function(fmt, ...)
    print(string.format(fmt, ...))
    pipe_out.Write(string.format(fmt, ...))
  end,
  trig = trig,
  pipe_out = pipe_out,
  pipe_in = pipe_in

}

return dsl
