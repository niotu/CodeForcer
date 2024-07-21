import React, {useState} from 'react';
import infoIcon from '../../assets/info-icon.png';

const InfoComponent = ({infoData}) => {
    const [showTooltip, setShowTooltip] = useState(false);

    const toggleTooltip = () => {
        setShowTooltip(!showTooltip); // Переключаем состояние showTooltip при клике
    };

    return (
        <div className="info-container">
            <span className="info-icon" onClick={toggleTooltip}><img src={infoIcon} alt={'i'} width={20} height={20}/></span>
            {showTooltip && (
                <div className="info-tooltip">
                    <p dangerouslySetInnerHTML={{__html: infoData.content}}></p>
                </div>
            )}
        </div>
    );
};

export default InfoComponent;
