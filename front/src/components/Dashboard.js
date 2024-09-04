import UserPage from "../pages/UserPage";
import BarChartCard from "./BarChartCard";
import LineChartCard from "./LineChartCard";
import PieChartCard from "./PieChartCard";
import SmallCard from "./SmallCard"
import HeaderTableCard from "./HeaderTableCard"
import { UserSidebar } from "./Sidebar";
import { useRef } from "react";

export default function Dashboard() {
    const ref1 = useRef(null)
    const ref2 = useRef(null)
    const ref3 = useRef(null)
    const ref4 = useRef(null)

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
                    <SmallCard title={"Total Requests"} number={1912} measurement={"req"} reff={ref1}/>
                    <SmallCard title={"Average Latency"} number={0.42} measurement={"s"}/>
                    <SmallCard title={"Success Rate"} number={97.3} measurement={"%"}/>
                    <SmallCard title={"Error Rate"} number={2.7} measurement={"%"}/>
                    <SmallCard title={"Requests Per Second"} number={1.2} measurement={"hits"}/>
                    <LineChartCard 
                        title={"Requests Per Day"}
                        x={[
                            "01/01/2000", "02/01/2000", "03/01/2000", "04/01/2000", "05/01/2000", "06/01/2000", 
                            "07/01/2000", "08/01/2000", "09/01/2000", "10/01/2000", "11/01/2000", "12/01/2000", 
                        ]}
                        y={[
                            302, 201, 189, 304, 190, 220,
                            250, 185, 203, 284, 196, 244
                        ]}
                        width={"64%"}
                        type="line"
                        height={"300px"}
                        reff={ref2}
                    />
                    <PieChartCard 
                        title={"Status Code Breakdown"}
                        x={[73, 4, 17, 6]}
                        y={["2XX", "3XX", "4XX", "5XX"]}
                        colors={["#222222", "#888888", "#999999", "#AAAAAA"]}
                        width={"35%"}
                        type="pie"
                        height={"300px"}
                    />
                    <PieChartCard 
                        title={"Request Method Breakdown"}
                        x={[45, 33, 4, 12, 6]}
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
                        y={[
                            302, 201, 189, 304, 190, 220
                        ]}
                        width={"64%"}
                        type="line"
                        height={"300px"}
                    />
                    <HeaderTableCard 
                        headersMap={requestHeaders} 
                        height={"300px"}
                        width={"49%"}
                        title={"Request Headers"}
                        reff={ref3}
                    />
                    <HeaderTableCard 
                        headersMap={responseHeaders} 
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