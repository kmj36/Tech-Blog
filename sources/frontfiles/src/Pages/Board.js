import Body from '../Components/parts/body';
import More from '../Components/parts/more';

function Board() {
    return(
        <>
        {Body(More())}
        </>
    );
}

export default Board;