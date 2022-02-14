package InitRedis

const SpikeCourseScript = `
	-- KEYS[1]: StudentID
    -- KEYS[2]: CourseID   
    -- 返回值有-1, -2, -3, 都代表抢课失败
    -- 返回值为1代表抢课成功
	-- -1代表没有这样的课程，-2代表课程容量已满，-3代表已经抢到课程

    -- 检查课程是否存在
	local nameOfSet = "course"..KEYS[2]
	local courseExists = redis.call("exists", KEYS[2]);
	if (courseExists == 0)
	then
		return -1;  -- 课程不存在
	end

	if (tonumber(redis.call("get", KEYS[2])) <= 0)  --- 课程容量已满
    then
		return -2; 
	end

    -- 检查当前用户是否已经选过当前课程 --
	local userHasCourse = redis.call("SISMEMBER", nameOfSet, KEYS[1]);
	if (userHasCourse == 1)
	then
		return -3;
	end

    -- 选课成功 --
	redis.call("decr", KEYS[2]);
	redis.call("sadd", nameOfSet, KEYS[1]);
	return 1;
`

var SpikeCourseSHA string

func init() {
	// 让redis加载抢课的lua脚本
	//
	SpikeCourseSHA = PrepareScript(SpikeCourseScript)

}
