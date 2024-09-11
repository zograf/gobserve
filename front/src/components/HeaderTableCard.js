import React from 'react';

export default function HeaderTableCard({headersMap, height, width, title, reff}) {
  const headersArray = Object.entries(headersMap);

  return (
    <div className='card' style={{height: height, width: width}} ref={reff}>
      <p>{title}</p>
      <div style={{overflowY: "scroll", height: '95%'}}>
        <table style={{ height: '100%', width: '100%', borderCollapse: 'collapse', marginTop: '20px', }}>
          <thead>
            <tr>
              <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd', width: "50%" }}>Header</th>
              <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd', width: "50%"}}>Value</th>
            </tr>
          </thead>
          <tbody>
            {headersArray.map(([header, value], index) => (
              <tr key={index}>
                <td style={{ padding: '8px', borderBottom: '1px solid #ddd' }}>{header}</td>
                <td style={{ padding: '8px', borderBottom: '1px solid #ddd' }}>{value}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};