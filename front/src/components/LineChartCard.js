import React, { useEffect, useState } from "react";
import Chart from "react-apexcharts";

export default function LineChartCard({title, x, y, width, height, reff}) {
    const [data, setData] = useState({
      options: {
        chart: {
          id: "basic-line",
          toolbar: {
            show: false
          }
        },
        xaxis: {
          categories: x
        },
        colors: ["#000000", "#c88214"]
      },
      series: [{
          data: y
        }]
    })

    useEffect(() => {
      setData(prevData => ({
        options: {
          chart: prevData.options.chart,
          xaxis: {
            categories: x
          },
          colors: prevData.options.colors
        },
        series: [{
          data: y
        }]
    }));
    }, [x, y])

    return (
        <div className="card" style={{width: width, height: height}} ref={reff}>
            <p className="chart-title">{title}</p>
            <Chart
                options={data.options}
                series={data.series}
                type="line"
                height={"100%"}
                width={"100%"}
            />
        </div>
    );
}