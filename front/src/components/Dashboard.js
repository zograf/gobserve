import UserPage from "../pages/UserPage";
import LineChartCard from "./LineChartCard";
import PieChartCard from "./PieChartCard";
import SmallCard from "./SmallCard"

export default function Dashboard() {

    return(
        <UserPage>
            <main className="mh-100">
                <h1 className="page-title">Welcome Back!</h1>
                <div className="flex wrap space-between gap-s">
                    <SmallCard title={"Total Requests"} number={1912} measurement={"req"}/>
                    <SmallCard title={"Average Latency"} number={0.42} measurement={"s"}/>
                    <SmallCard title={"Success Rate"} number={97.3} measurement={"%"}/>
                    <SmallCard title={"Error Rate"} number={2.7} measurement={"%"}/>
                    <SmallCard title={"Requests Per Second"} number={1.2} measurement={"hits"}/>
                    <LineChartCard 
                        title={"Requests per day"}
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
                    />
                    <PieChartCard 
                        title={"Status Codes"}
                        x={[73, 4, 17, 6]}
                        y={["2XX", "3XX", "4XX", "5XX"]}
                        colors={["#222222", "#888888", "#999999", "#AAAAAA"]}
                        width={"35%"}
                        type="pie"
                        height={"300px"}
                    />
                </div>
            </main>
        </UserPage>
    )
}