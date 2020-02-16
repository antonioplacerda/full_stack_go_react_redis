import React from 'react';
import './App.css';

import Jobs from './Jobs';

const JOBS_API_URL = 'http://localhost:8080/jobs';

async function fecthJobs(updateCB) {
  const res = await fetch(JOBS_API_URL);
  let json = await res.json();

  updateCB(json);
}

function App() {

  const [jobList, updateJobs] = React.useState([]);

  React.useEffect(() => {
    fecthJobs(updateJobs);
  }, [])

  return (
    <div className="App">
      <Jobs jobs={jobList}/>
    </div>
  );
}

export default App;
