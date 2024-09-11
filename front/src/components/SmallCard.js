import "./SmallCard.css"

export default function SmallCard({title, number, measurement, reff}) {

    return(
        <div className="card" ref={reff}>
            <p className="smallcard-title">{title}</p>
            <div >
                <p className="smallcard-number">{number}</p>
                <p className="smallcard-measurement">{measurement}</p>
            </div> 
        </div>
    )
}