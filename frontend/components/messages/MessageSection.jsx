import React, { Component } from 'react'
import PropTypes from 'prop-types'
import MessageList from './MessageList.jsx'
import MessageForm from './MessageForm.jsx'

class MessageSection extends Component{

    render(){
        let {activeChannel} = this.props
        console.log(this.props.messages)
        return(
            <div className='messages-container panel panel-default'>
                <div className='panel-heading'>
                    <strong>{activeChannel !==undefined? activeChannel.name: 'Set the Channel'}</strong>
                </div>
                <div className='panel-body messages'>
                    <MessageList {...this.props} />
                    <MessageForm {...this.props} />
                </div>
            </div>
        )
    }
}
MessageSection.propTypes = {
    messages : PropTypes.array.isRequired,
    addMessage : PropTypes.func.isRequired,
    activeChannel :PropTypes.object.isRequired,
}

export default MessageSection
