package routes

import (
	"delay_queue_callback/core"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func Handle(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		get(&response, request)
	case "POST":
		add(&response, request)
	case "DELETE":
		del(&response, request)
	default:
		response.Write(createResponseBody(0, "不存在的请求方法", nil))
	}

	return
}
func add(response *http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		(*response).Write(createResponseBody(0, "解析body失败", nil))
		return
	}

	var job core.Job
	err = json.Unmarshal(body, &job)

	if err != nil {
		(*response).Write(createResponseBody(0, err.Error(), nil))
		return
	}
	if job.Id == "" {
		(*response).Write(createResponseBody(0, "不合法的JOBID", nil))
		return
	}

	if job.Delay <= 0 {
		(*response).Write(createResponseBody(0, "执行时间不合法", nil))
		return
	}

	if job.Callback == "" {
		(*response).Write(createResponseBody(0, "未填写回调URL", nil))
		return
	}

	if job.Topic == "" {
		(*response).Write(createResponseBody(0, "工作类型不能为空", nil))
		return
	}

	job.JobSign = core.JobSign(job.Topic, job.Id)

	err = core.NewWork(job)
	if err != nil {
		(*response).Write(createResponseBody(0, "加入工作队列失败", nil))
		return
	}

	(*response).Write(createResponseBody(1, "添加新任务成功", nil))

}

func get(response *http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()

	if err != nil {
		(*response).Write(createResponseBody(0, "解析命令行参数失败", nil))
		return
	}
	jobId := request.Form.Get("id");
	if strings.TrimSpace(jobId) == "" {
		(*response).Write(createResponseBody(0, "jobID参数缺失", nil))
		return
	}
	topic := request.Form.Get("topic");
	if strings.TrimSpace(topic) == "" {
		(*response).Write(createResponseBody(0, "topic参数缺失", nil))
		return
	}
	job, err := core.GetJob(core.JobSign(topic, jobId))

	if err != nil {
		(*response).Write(createResponseBody(0, "获取job信息失败", nil))
		return

	}

	(*response).Write(createResponseBody(1, "", job))
	return

}

func del(response *http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()

	if err != nil {
		(*response).Write(createResponseBody(0, "解析命令行参数失败", nil))
		return
	}

	jobId := request.Form.Get("id");
	if strings.TrimSpace(jobId) == ""{
		(*response).Write(createResponseBody(0, "jobId不能为空", nil))
		return
	}
	topic := request.Form.Get("topic");
	if strings.TrimSpace(topic) == ""{
		(*response).Write(createResponseBody(0, "topic参数不能为空", nil))
		return
	}
	err = core.RemoveJob(core.JobSign(topic, jobId));
	if err != nil {
		(*response).Write(createResponseBody(0, "删除job失败", nil))
		return
	}

	(*response).Write(createResponseBody(1, "删除job成功", nil))

}

type Response struct {
	Success int         `json:"success"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func createResponseBody(success int, msg string, data interface{}) ([]byte) {
	response := &Response{success, msg, data}
	val, err := json.Marshal(response)

	if err != nil {
		return []byte(`{"success":"0", "message": "解析返回数据失败", "data":null}`)
	}

	return val

}
