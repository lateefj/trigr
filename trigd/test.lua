print("Hello triggerd stin:\n")

print(string.format("Hello from %s trigger type %s\n",_VERSION, trigr.Type))
function handle_log(log)
  print(string.format("Log Text: %s", log.Text))
end
print(trigr.Data)
