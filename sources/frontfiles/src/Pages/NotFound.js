import Body from '../Components/parts/body';
import styled from 'styled-components';

const Flexbox = styled.div`
    display: flex;
    flex-direction: column;
    width: 100%;
    height: 100%;
    justify-content: center;
    align-items: center;
`;

function notfound(){
    return(
        Body(
            <Flexbox>
                <img src="404.png" alt="404" width="400px"/>
                <h1>
                    404 NotFound
                </h1>
            </Flexbox>
        )
    );
}

export default notfound; 