import groovy.json.JsonSlurper
import java.util.logging.Logger
import java.util.logging.Level


def fetchEnvData() {
    def mapping = []
    def envUrl = System.getenv("ENV_DATA_URL") ?: "http://127.0.0.1:9301/env"
    def logger = Logger.getLogger("env.fetcher")

    try {
        logger.info("EnvRequest: 准备请求接口 -> ${envUrl}")

        def connection = new URL(envUrl).openConnection()
        connection.setConnectTimeout(2000)
        connection.setReadTimeout(5000)
        connection.setRequestProperty("Accept", "application/json")

        def responseCode = connection.responseCode
        if (responseCode != 200) {
            logger.warning("EnvRequest: 接口响应异常，HTTP Code: ${responseCode}")
            return ["接口响应异常 (HTTP ${responseCode})"]
        }

        def jsonRaw = connection.inputStream.text
        def jsonObj = new JsonSlurper().parseText(jsonRaw)
        def lastMachine = ""
        mapping.add(" ")
        jsonObj.each { item ->
            def currentMachine = item.ident?.split(':')?.getAt(0)?.split(/\./)?.getAt(0) ?: ""

            if (lastMachine && lastMachine != currentMachine) {
                mapping.add("-")
            }
            mapping.add("${item.ident} | ${item.date ?: 'N/A'} | ${item.owner ?: 'N/A'}")
            lastMachine = currentMachine
        }

        logger.info("EnvRequest: 处理完成，生成选项数量: ${mapping.size()}, 最终 mapping 列表内容: ${mapping.inspect()}")
    } catch (Exception e) {
        logger.log(Level.SEVERE, "EnvRequest: 脚本运行异常 -> ${e.message}", e)
        return ["执行失败，请检查 Jenkins 系统日志"]
    }
    return mapping
}

return fetchEnvData()