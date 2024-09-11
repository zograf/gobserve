import UserPage from "../pages/UserPage";
import BarChartCard from "./BarChartCard";
import LineChartCard from "./LineChartCard";
import PieChartCard from "./PieChartCard";
import SmallCard from "./SmallCard"
import HeaderTableCard from "./HeaderTableCard"
import { UserSidebar } from "./Sidebar";
import { useEffect, useRef, useState } from "react";
import axios from "axios";

export default function Dashboard() {
    const ref1 = useRef(null)
    const ref2 = useRef(null)
    const ref3 = useRef(null)
    const ref4 = useRef(null)
    const [data, setData] = useState([])
    const [totalRequests, setTotalRequests] = useState(0)
    const [averageLatency, setAverageLatency] = useState(0)
    const [successRate, setSuccessRate] = useState(0)
    const [errorRate, setErrorRate] = useState(0)
    const [requestsPerSecond, setRequestsPerSecond] = useState(0.0)
    const [statusSummary, setStatusSummary] = useState([0, 0, 0, 0])
    const [dates, setDates] = useState([])
    const [dateValues, setDateValues] = useState([])
    const [responseTimeBar, setResponseTimeBar] = useState([])
    const [methods, setMethods] = useState([])
    const [reqHead, setReqHead] = useState({})
    const [resHead, setResHead] = useState({})

    const requestHeaders = {
        "Accept": "image/avif",
        "Accept-Encoding": "gzip",
        "Accept-Language": "en-US",
        "Connection": "keep-alive",
        "Referer": "http://localhost:42547/log",
        "Sec-Ch-Ua": "\"Not)A;Brand\";v=\"99\"",
        "Sec-Ch-Ua-Mobile": "?0",
        "Sec-Ch-Ua-Platform": "\"Linux\"",
        "Sec-Fetch-Dest": "image",
        "Sec-Fetch-Mode": "no-cors",
        "Sec-Fetch-Site": "same-origin",
        "Sec-Gpc": "1",
        "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36"
      };

    const responseHeaders = { 
        "Content-Type": "application/json", 
        "Content-Length": "123", 
        "Cache-Control": "no-cache" 
    };

    useEffect(() => {
        axios.get("http://localhost:8080/agg/log")
          .then((response) => {
            let d = response.data.Data
            setData(d)
            calculateNumbers(d)
          })
          .catch((err) => { 
            console.log(err)
          });
      }, []);
    
      const calculateNumbers = (data) => {
        let totalRequests = 0
        let totalLatency = 0
        let successfulRequests = 0
        let minTimestamp = -Infinity
        let maxTimestamp = Infinity
        let sSummary = [0, 0, 0, 0]
        let dateCount = {}
        let bins = {
            '0-10ms': 0,
            '10-20ms': 0,
            '20-30ms': 0,
            '30-40ms': 0,
            '40-50ms': 0,
            '50+ms': 0
        };
        let methodCount = {
            "GET": 0,
            "POST": 0,
            "PUT": 0,
            "DELETE": 0,
            "OPTIONS": 0
        }
        
        let reqHeaderSet = new Set()
        let resHeaderSet = new Set()

        Object.keys(data).forEach(ms => {
            data[ms].forEach(request => {
                const duration = request.duration_ms
                totalLatency += request.duration_ms;
                totalRequests++; 
                const statusCode = request.status_code
                const method = request.method
                const reqHeaders = request.request_headers
                const resHeaders = request.response_headers

                if (statusCode >= 200 && statusCode < 300) {
                    successfulRequests++; 
                }
                if (statusCode >= 200 && statusCode < 300) {
                    sSummary[0]++;
                } else if (statusCode >= 300 && statusCode < 400) {
                    sSummary[1]++;
                } else if (statusCode >= 400 && statusCode < 500) {
                    sSummary[2]++;
                } else if (statusCode >= 500 && statusCode < 600) {
                    sSummary[3]++;
                }

                const requestTime = new Date(request.request_timestamp).getTime();
                if (requestTime < minTimestamp) {
                    minTimestamp = requestTime;
                }
                if (requestTime > maxTimestamp) {
                    maxTimestamp = requestTime;
                }

                const date = request.request_timestamp.split('T')[0];
                if (dateCount[date]) {
                    dateCount[date]++;
                } else {
                    dateCount[date] = 1;
                }

                if (duration >= 0 && duration < 10) {
                    bins["0-10ms"]++
                } else if (duration > 10 && duration < 20) {
                    bins["10-20ms"]++
                } else if (duration > 20 && duration < 30) {
                    bins["20-30ms"]++
                } else if (duration > 30 && duration < 40) {
                    bins["30-40ms"]++
                } else if (duration > 40 && duration < 50) {
                    bins["40-50ms"]++
                } else {
                    bins["50+ms"]++
                }
                methodCount[method]++
                for (const [key, value] of Object.entries(reqHeaders)) {
                    reqHeaderSet.add(`${key}: ${value}`);
                }
                for (const [key, value] of Object.entries(resHeaders)) {
                    resHeaderSet.add(`${key}: ${value}`);
                }
            });
          });
        let timeSpanSeconds = (maxTimestamp - minTimestamp) / 1000;
        let rps = totalRequests / timeSpanSeconds;

        setTotalRequests(totalRequests)
        setAverageLatency(totalLatency/totalRequests)
        setSuccessRate((successfulRequests / totalRequests) * 100)
        setErrorRate(((totalRequests - successfulRequests) / totalRequests) * 100)
        setRequestsPerSecond((rps+0.01).toFixed(2))
        setStatusSummary(sSummary)
        setDates(Object.keys(dateCount))
        setDateValues(Object.values(dateCount))
        setResponseTimeBar(Object.values(bins))
        setMethods(Object.values(methodCount))

        let reqHead = {}
        reqHeaderSet.forEach((entry) => {
            reqHead[entry.split(": ")[0]] = entry.split(": ")[1]
        });

        let resHead = {}
        resHeaderSet.forEach((entry) => {
            resHead[entry.split(": ")[0]] = entry.split(": ")[1]
        });

        setReqHead(reqHead)
        setResHead(resHead)
      }


    return(
        <main className="mh-100">
            <div className="sidebar-root">
                <UserSidebar ref1={ref1} ref2={ref2} ref3={ref3} ref4={ref4}/>
                <div className="header-root">
                    <div>

                    </div>
                    <div className="page-wrapper">
            <main className="mh-100" ref={ref4}>
                <h1 className="page-title">Welcome Back!</h1>
                <div className="flex wrap space-between gap-s">
                    <SmallCard title={"Total Requests"} number={totalRequests} measurement={"req"} reff={ref1}/>
                    <SmallCard title={"Average Latency"} number={averageLatency} measurement={"ms"}/>
                    <SmallCard title={"Success Rate"} number={successRate} measurement={"%"}/>
                    <SmallCard title={"Error Rate"} number={errorRate} measurement={"%"}/>
                    <SmallCard title={"Requests Per Second"} number={requestsPerSecond} measurement={"hits"}/>
                    <LineChartCard 
                        title={"Requests Per Day"}
                        x={dates}
                        y={dateValues}
                        width={"64%"}
                        type="line"
                        height={"300px"}
                        reff={ref2}
                    />
                    <PieChartCard 
                        title={"Status Code Breakdown"}
                        x={statusSummary}
                        y={["2XX", "3XX", "4XX", "5XX"]}
                        colors={["#222222", "#888888", "#999999", "#AAAAAA"]}
                        width={"35%"}
                        type="pie"
                        height={"300px"}
                    />
                    <PieChartCard 
                        title={"Request Method Breakdown"}
                        x={methods}
                        y={["GET", "POST", "PUT", "DELETE", "OPTION"]}
                        colors={["#222222", "#888888", "#999999", "#AAAAAA", "#CCCCCC"]}
                        width={"35%"}
                        type="pie"
                        height={"300px"}
                    />
                    <BarChartCard 
                        title={"Response Time Distribution"}
                        x={[
                            "0ms-5ms", "5ms-10ms", "10ms-15ms", "15ms-20ms", "20ms-25ms", "25ms-30ms"
                        ]}
                        y={responseTimeBar}
                        width={"64%"}
                        type="line"
                        height={"300px"}
                    />
                    <HeaderTableCard 
                        headersMap={reqHead} 
                        height={"300px"}
                        width={"49%"}
                        title={"Request Headers"}
                        reff={ref3}
                    />
                    <HeaderTableCard 
                        headersMap={resHead} 
                        height={"300px"}
                        width={"49%"}
                        title={"Response Headers"}
                    />
                </div>
                <br></br>
            </main>
                    </div>
                </div>
            </div>
        </main>
    )
}