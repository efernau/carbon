-- MessagePack Wrappers

-- Load 3rdparty libraries
-- MessagePack.lua
local msgpack = assert(loadstring(carbon.glue("3rdparty/MessagePack.lua")))()

-- Support more lua types, but at the cost of compatibility with non-carbon msgpack things.
-- Functions
msgpack.packers['_function'] = function(buffer, fn)
	return msgpack.packers['ext'](buffer, 7, assert(string.dump(fn)))
end
msgpack.packers['function'] = function(buffer, fn)
	local upvals = debug.getallupvals(fn)
	if upvals then
		local buf = {}
		msgpack.packers['_function'](buf, fn)
		msgpack.packers['table'](buf, upvals)
		msgpack.packers['ext'](buffer, 8, table.concat(buf))
	else
		msgpack.packers['_function'](buffer, fn)
	end
end

-- Tables
msgpack.packers['table'] = function(buffer, t)
	local mt = getmetatable(t)
	if mt then
		local buf = {}
		msgpack.packers['_table'](buf, t)
		msgpack.packers['table'](buf, mt)
		msgpack.packers['ext'](buffer, 42, table.concat(buf))
	else
		msgpack.packers['_table'](buffer, t)
	end
end

-- Userdata
msgpack.packers['userdata'] = function(buffer, ud)
	-- Hack. A big one, too.
	local idpool = "tmp:msgpack:ud-tmp"
	kvstore.inc(idpool)
	local tmpid = tostring(kvstore.get(idpool))..":"..carbon.randomstring(16)
	local id = idpool..":"..tmpid
	kvstore.set(id, ud)
	msgpack.packers['ext'](buffer, 43, tmpid)
end

-- Unpacker for both
msgpack.build_ext = function (tag, data)
	if tag == 7 then -- Function
		return assert(loadstring(data))
	elseif tag == 8 then -- Function with upvalues
		local f = msgpack.unpacker(data)
		local _, fn = f()
		local _, upvals = f()
		debug.setallupvals(fn, upvals)
		return fn
	elseif tag == 42 then -- Table
		local f = msgpack.unpacker(data)
		local _, t = f()
		local _, mt = f()
		return setmetatable(t, mt)
	elseif tag == 43 then -- Userdata. Big hack.
		local id = "tmp:msgpack:ud-tmp:"..data
		local ud = kvstore.get(id)
		kvstore.del(id)
		return ud
	end
end

return msgpack
