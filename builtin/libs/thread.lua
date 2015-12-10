-- Threads
-- Has a buttload of issues, mostly bindings. Gotta looooove broken shit.

local msgpack = require("msgpack")

function thread.spawn(fn, bindings, buffer)
	local code = ""
	if type(fn) == "function" then
		code = string.dump(fn)
	elseif type(fn) == "string" then
		fn, err = loadstring(code)
		if not err then
			code = string.dump(fn)
		else
			error(err)
		end
	end

	local buffer = tonumber(buffer) or -1

	local ch
	local err
	if type(bindings) == "table" then
		ch, err = thread._spawn(code, true, bindings, buffer)
	else
		ch, err = thread._spawn(code, false, {["s"]="v"}, buffer)
	end
	if err ~= nil then
		return false, error(err)
	end
	return ch
end

function thread.rpcthread() -- not working, issues with binding or something .-.
	local rpc = thread.spawn(function()
	local msgpack = require("msgpack")
	while true do
		local src = com.receive(threadcom)
		local args = msgpack.unpack(src)
		local f, err = loadstring(args.f)
		if not err then
			com.send(threadcom, msgpack.pack({pcall(f, unpack(args.args))}))
		else
			com.send(threadcom, msgpack.pack({false, err}))
		end
	end
end)

function call(f, ...)
	com.send(rpc, msgpack.pack({
		f = f,
		args = {...}
	}))
end

	local function recieve()
		return unpack(msgpack.unpack(com.receive(rpc)))
	end

	return {
		["call_async"] = call,
		["call"] = (function(f, ...)
			call(f, ...)
			return recieve()
		end),
		["recieve"] = recieve,
	}
end

function thread.kvstore() -- doesn't work either .-.
	local chan = thread.spawn(function()
		local msgpack = require("msgpack")
		local store = {}
		while true do
			local suc, cmd = pcall(msgpack.unpack, com.receive(threadcom))
			if suc then
				if cmd.value then
					store[cmd.name] = cmd.value
					com.send(threadcom, msgpack.pack({value=true, error=nil}))
				else
					com.send(threadcom, msgpack.pack({value=store[cmd.name], error=nil}))
				end
			else
				com.send(threadcom, msgpack.pack({value=nil, error=cmd}))
			end
		end
	end)

	return function(name, value)
		if name then
			com.send(chan, msgpack.pack({["name"]=name, ["value"]=value}))
			local res = com.receive(chan)
			if res.error then
				return false, res.error
			else
				return true, res.value
			end
		end
	end
end

return thread