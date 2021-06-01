import React from 'react';
import './body.css';

function Body(value) {
    return(
        <div className = "middler_con">
            <div className = "midwrapper">
                {value}
            </div>
        </div>  
    );
}

export default Body;