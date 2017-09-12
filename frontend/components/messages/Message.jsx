import React, { Component } from 'react'
import PropTypes from 'prop-types'
import fecha from 'fecha'

class Message extends Component{

    render(){
        let {message} = this.props;
        
        let createdAt = fecha.format(new Date(message.createdAt),'YYYY-MM-DD hh:mm A')
        console.log(message)
    

        return(
            <li className='message'>
                <div className='author'><strong>{message.Author}</strong><i className='timestamp'>{createdAt}</i></div>
                <div className='body'>{message.body}</div>
            </li>
        )
    }
}
Message.propTypes={
    message: PropTypes.object.isRequired,
}
export default Message