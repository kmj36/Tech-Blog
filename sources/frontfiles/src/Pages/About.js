import Body from '../Components/parts/body';
import './About.css';
function About() {
    document.title = "About Me"
    return(
        Body(
            <div className = "Aboutcontainer">
                <img className = "aboutlogo" src="K.svg"/>
                <div className = "Aboutcolume">
                    <div><p>My Name Is </p><h2>Kim Min-Je, Network Security Software Developer</h2></div>
                    <div>
                        <h3 className = "contact">Contact</h3>
                        <a href = "https://github.com/kmj36"><img className="snsicon" src="githubicon.png"/></a>
                        <a href = "mailto:kmj36953695@gmail.com"><img className="snsicon" src="emailicon.jpg"/></a>
                        <footer>
                            <p>Copyright 2021. Kimminje All pictures cannot be copied without permission.</p>
                        </footer>
                    </div>
                </div>
            </div>
        )
    );
}

export default About;