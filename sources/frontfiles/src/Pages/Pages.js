import Body from '../Components/parts/body';
import GetPostandComm from '../Components/parts/getpostandcomment';
import qs from 'query-string';

function pages(location) {
    const id = qs.parse(location.location.search)
    return(
        <>
        {Body(GetPostandComm(id.id))}
        </>
    );
};


export default pages;