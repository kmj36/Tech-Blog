import Body from '../Components/parts/body';
import Uploadform from '../Components/forms/uploadform';
import qs from 'query-string';

function Upload(location) {
    const id = qs.parse(location.location.search)
    return(
        <>
        {Body(Uploadform(id.load))}
        </>
    );
}

export default Upload;