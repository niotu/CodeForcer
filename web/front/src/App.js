import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import LoginPage from './components/LoginPage';
import GroupsPage from './components/GroupsPage';
import ContestsPage from './components/ContestsPage';
import ContestDetails from './components/ContestDetails';

const App = () => {
    return (
        <Router>
            <Routes>
                <Route path="/" element={<LoginPage />} />
                <Route path="/groups" element={<GroupsPage />} />
                <Route path="/contests/:groupCode" element={<ContestsPage />} />
                <Route path="/contest-details/:groupCode/:contestId" element={<ContestDetails />} />
            </Routes>
        </Router>
    );
};

export default App;
